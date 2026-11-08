// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package game

import "time"

const stageTransitionDuration = 2 * time.Second

type stageTransition struct {
	elapsed       time.Duration
	startSpeed    float64
	targetSpeed   float64
	weatherStartX float64
	weatherEndX   float64
	removeWeather bool
}

func (m *Model) inheritWorld(previous *Model) {
	for _, terrain := range m.terrain {
		for _, chunk := range terrain.chunks {
			m.scene.Remove(chunk)
		}
	}
	m.terrain = previous.terrain
	for _, terrain := range m.terrain {
		for _, chunk := range terrain.chunks {
			m.scene.Add(chunk)
		}
	}

	if previous.stage.DarkCloudsOnly || m.stage.DarkCloudsOnly {
		for _, cloud := range previous.clouds {
			if _, visible := cloud.sprite.ActiveFrame(); visible {
				m.scene.Add(cloud.sprite)
				m.departing = append(m.departing, cloud.sprite)
			}
		}
		for _, storm := range previous.storms {
			if storm == nil {
				continue
			}
			storm.bolt.Hide()
			m.scene.Add(storm.cloud)
			m.departing = append(m.departing, storm.cloud)
		}
		for _, hazard := range []*pterosaurHazard{previous.pterosaur, previous.redEye} {
			if hazard != nil {
				m.scene.Add(hazard.sprite)
				m.departing = append(m.departing, hazard.sprite)
			}
		}
	}
	for _, item := range previous.departing {
		m.scene.Add(item)
		m.departing = append(m.departing, item)
	}
}

func (m *Model) beginStageTransition(previous *Model) {
	m.scrollSpeed = previous.scrollSpeed
	m.transition = &stageTransition{
		startSpeed:  previous.scrollSpeed,
		targetSpeed: m.stage.ScrollSpeed,
	}
	m.applyScrollSpeed(m.scrollSpeed)

	if m.stage.DarkCloudsOnly {
		for index, storm := range m.storms {
			storm.setPosition(float64(m.Width+10+index*24), m.cloudY(storm.yFraction, 5))
		}
	}
	if previous.stage.DarkCloudsOnly && !m.stage.DarkCloudsOnly {
		for index, cloud := range m.clouds {
			_, y := cloud.sprite.Position()
			cloud.sprite.SetPosition(float64(m.Width+10+index*18), y)
		}
	}

	switch {
	case !previous.stage.PersistentRain && m.stage.PersistentRain:
		m.transition.weatherStartX = float64(m.Width)
		m.transition.weatherEndX = 0
		m.positionWeather(m.transition.weatherStartX)
	case previous.stage.PersistentRain && !m.stage.PersistentRain:
		if m.rain != nil {
			m.scene.Remove(m.rain)
		}
		if m.ceiling != nil {
			m.scene.Remove(m.ceiling.cloud)
			m.scene.Remove(m.ceiling.bolt)
		}
		m.rain = previous.rain
		m.ceiling = previous.ceiling
		if m.rain != nil {
			m.scene.Add(m.rain)
		}
		if m.ceiling != nil {
			m.scene.Add(m.ceiling.cloud)
			m.scene.Add(m.ceiling.bolt)
		}
		m.transition.weatherStartX = 0
		m.transition.weatherEndX = -float64(m.Width)
		m.transition.removeWeather = true
		m.positionWeather(0)
	}
}

func (m *Model) updateStageTransition(step time.Duration) {
	if m.transition == nil {
		return
	}
	m.transition.elapsed = min(m.transition.elapsed+step, stageTransitionDuration)
	progress := float64(m.transition.elapsed) / float64(stageTransitionDuration)
	// Smoothstep keeps both ends of the wind change from feeling abrupt.
	progress = progress * progress * (3 - 2*progress)
	m.scrollSpeed = m.transition.startSpeed + (m.transition.targetSpeed-m.transition.startSpeed)*progress
	m.applyScrollSpeed(m.scrollSpeed)
	weatherX := m.transition.weatherStartX + (m.transition.weatherEndX-m.transition.weatherStartX)*progress
	m.positionWeather(weatherX)
	if m.transition.elapsed < stageTransitionDuration {
		return
	}
	if m.transition.removeWeather {
		if m.rain != nil {
			m.scene.Remove(m.rain)
			m.rain = nil
		}
		if m.ceiling != nil {
			m.scene.Remove(m.ceiling.cloud)
			m.scene.Remove(m.ceiling.bolt)
			m.ceiling = nil
		}
	}
	m.scrollSpeed = m.transition.targetSpeed
	m.applyScrollSpeed(m.scrollSpeed)
	m.transition = nil
}

func (m *Model) positionWeather(x float64) {
	if m.rain != nil {
		m.rain.SetPosition(x, 0)
	}
	if m.ceiling != nil {
		m.ceiling.setOffset(x)
	}
}

func (m *Model) applyScrollSpeed(speed float64) {
	for index, terrain := range m.terrain {
		terrainSpeed := -speed
		if index != len(m.terrain)-1 {
			terrainSpeed = -0.55 * speed / worldScrollSpeed
		}
		terrain.speed = terrainSpeed
		for _, chunk := range terrain.chunks {
			chunk.SetVelocity(terrainSpeed, 0)
		}
	}
	for _, tree := range m.trees {
		tree.sprite.SetVelocity(-speed, 0)
	}
	for _, fern := range m.ferns {
		fern.sprite.SetVelocity(-speed, 0)
	}
	for _, grazer := range m.grazers {
		grazer.sprite.SetVelocity(-speed, 0)
	}
	if m.finish != nil && !m.completed {
		m.finish.SetVelocity(-speed, 0)
	}
	for _, storm := range m.storms {
		if storm == nil {
			continue
		}
		velocity := stormSpeed * speed / worldScrollSpeed
		storm.cloud.SetVelocity(velocity, 0)
		storm.bolt.SetVelocity(velocity, 0)
	}
}

func (m *Model) updateDepartingScenery() {
	kept := m.departing[:0]
	for _, item := range m.departing {
		frame, visible := item.ActiveFrame()
		if !visible {
			m.scene.Remove(item)
			continue
		}
		x, _ := item.Position()
		if x+float64(frame.Width) < 0 {
			m.scene.Remove(item)
			continue
		}
		kept = append(kept, item)
	}
	m.departing = kept
}

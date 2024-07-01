# Scared sand design document

Traveling on a colossal holy land ship to make connection between surviving conclave

## Concept
- Survive and complete objective by racing with the DPS

## State machines
- Land ships (Idle, Dock, Cruising, Destroyed)
- Knights (Idle, Resting, Patrolling, In combat, Died)
- Land ship Components (Un-equipped, Idle, Active, Destroyed)
- Lost souls (Idle, Hunting, Destroyed)

## Maths
- Hit scan
- Path finding calculation
- Trajectory calculation

## Life cycles
- Momentum
- Path finding/Trajectory calculation
- Hit scan detection
- Spawning

## Stat properties
- HP
- AP (Flat damage reduction)
- Damage
- RPM (For projectile weapons, to determine number of hit)
- Range
- Pattern

## Jun 26
- When start a new expedition, select r number of rune templates
- During expedition, spend devotion to create runes
- Knight commands: To target, Move along, Guard Area
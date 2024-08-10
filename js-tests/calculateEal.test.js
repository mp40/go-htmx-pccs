import { expect, test } from 'vitest'
import {calculateEAL} from '../static/scripts/calculateEal'

test.each([
    {level: 1, range: 5, expected: 34},
    {level: 3, range: 5, expected: 38},
    {level: 5, range: 5, expected: 40},
    {level: 7, range: 5, expected: 42},
    {level: 1, range: 25, expected: 22},
    {level: 3, range: 25, expected: 26},
    {level: 5, range: 25, expected: 28},
    {level: 7, range: 25, expected: 30},
    {level: 1, range: 100, expected: 12},
    {level: 3, range: 100, expected: 16},
    {level: 5, range: 100, expected: 18},
    {level: 7, range: 100, expected: 20},
    {level: 1, range: 200, expected: 7},
    {level: 3, range: 200, expected: 11},
    {level: 5, range: 200, expected: 13},
    {level: 7, range: 200, expected: 15},
])('single shot by level $level stationary standing shooter at stationary target at $range', ({level, range, expected}) => {
    const result = calculateEAL(level, range, false, "Standing", "Standing Exposed", 0, 0)
    expect(result).toBe(expected)
})

test.each([
    {auto: false, target: "Standing Exposed", expected: 12},
    {auto: false, target: "Kneeling Exposed", expected: 11},
    {auto: false, target: "Fire Over/Around", expected: 5},
    {auto: false, target: "Look Over/Around", expected: 1},
    {auto: false, target: "Prone/Crawl", expected: 7},
    {auto: true, target: "Standing Exposed", expected: 19},
    {auto: true, target: "Kneeling Exposed", expected: 16},
    {auto: true, target: "Fire Over/Around", expected: 7},
    {auto: true, target: "Look Over/Around", expected: 2},
    {auto: true, target: "Prone/Crawl", expected: 7},
])('$auto auto shot by level 0 stationary shooter at stationary $target target at range 50', ({auto, target, expected}) => {
    const result = calculateEAL(0, 50, auto, "Standing", target, 0, 0)
    expect(result).toBe(expected)
})

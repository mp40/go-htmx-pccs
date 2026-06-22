import { expect, test } from 'vitest'
import { getBaseHitLocationSpacingFromSalm } from '../static/scripts/shotguns'

const testCases = [
  [-12, 1],
  [-11, 1],
  [-10, 2],
  [-7, 2],
  [-6, 3],
  [-5, 3],
  [-4, 4],
  [-3, 4],
  [-2, 6],
  [-1, 6],
  [0, 8],
  [1, 8],
  [2, 11],
  [3, 11],
  [4, 14],
  [5, 14],
  [6, 19],
  [7, 19],
  [8, 25],
  [9, 25],
  [10, 34],
  [11, 35],
  [12, 45],
  [13, 45],
  [14, 60],
  [15, 60],
  [16, 79],
  [17, 79],
  [18, 100],
]

test.each(testCases)("salm %i equals hit location spacing %i", (salm, hls) => {
  expect(getBaseHitLocationSpacingFromSalm(salm)).toBe(hls)
})
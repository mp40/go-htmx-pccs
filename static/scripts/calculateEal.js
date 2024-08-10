const rangeTable = Object.freeze({
  1: 33,
  2: 28,
  3: 25,
  4: 23,
  5: 22,
  6: 20,
  7: 19,
  8: 18,
  9: 17,
  11: 16,
  12: 15,
  14: 14,
  16: 13,
  19: 12,
  22: 11,
  25: 10,
  30: 9,
  35: 8,
  40: 7,
  45: 6,
  50: 5,
  55: 4,
  65: 3,
  75: 2,
  85: 1,
  100: 0,
  115: -1,
  130: -2,
  150: -3,
  170: -4,
  200: -5,
  230: -6,
  250: -7,
  300: -8,
  350: -9,
  400: -10,
  450: -11,
  500: -12,
  600: -13,
  700: -14,
});
  
const salTable = Object.freeze({
  0: 0,
  1: 5,
  2: 7,
  3: 9,
  4: 10,
  5: 11,
  6: 12,
  7: 13,
  8: 14,
  9: 15,
  10: 16,
  11: 17,
  12: 18,
  13: 19,
  14: 20,
  15: 21,
  16: 22,
  17: 23,
  18: 24,
  19: 25,
  20: 26,
});

const targetSizeMods = {
  "Look Over/Around": [-4, -3, -3],
  "Fire Over/Around": [0, 2, 2],
  "Standing Exposed": [7, 14, 1],
  "Kneeling Exposed": [6, 11, 3],
  "Prone/Crawl": [2, 2, 2],
  Running: [8, 14, 1],
  "Low Crouch": [7, 11, 2],
  "Hands & Knees": [6, 8, 1],
  "Low Prone": [1, 0, 5],
};

const shooterStanceMods = {
  Standing: 0,
  "Standing & Braced": 4,
  Kneeling: 3,
  "Kneeling & Braced": 5,
  Prone: 6,
  "Prone & Braced": 7,
};
  
// const movement = [0, 0.25, 0.5, 1, 1.5, 2, 2.5, 3];  
// HPI: [10, 20, 40, 70, 100]
const movementMods = {
  0: [0, 0, 0, 0, 0],
  0.5: [-6, -5, -5, -5, -5],
  1: [-8, -6, -5, -5, -5],
  2: [-10, -8, -6, -5, -5],
  3: [-10, -10, -7, -6, -5],
  4: [-10, -10, -8, -6, -6],
  10: [-10, -10, -10, -10, -8],
};
  
const getHpi = (combinedMove) => {
  if (combinedMove === 0) {
    return 0;
  }

  if (combinedMove <= 0.5) {
    return 0.5;
  }

  if (combinedMove > 0.5 && combinedMove <= 1) {
    return 1;
  }

  if (combinedMove > 1 && combinedMove <= 4) {
    return Math.ceil(combinedMove);
  }

  return 10;
};
  
const getRangeIndex = (range) => {
  if (range <= 10) {
    return 0;
  }

  if (range <= 20) {
    return 1;
  }

  if (range <= 40) {
    return 2;
  }

  if (range <= 70) {
    return 3;
  }

  return 4;
};

const getMovementMod = (shooterMove, targetMove, range) => {
  const move = shooterMove + targetMove;
  const hpi = getHpi(move);
  const index = getRangeIndex(range);

  return movementMods[hpi][index];
};

 /**
 * Calculates the shooting EAL based on level, fire mode, stance, size and movement.
 *
 * @param {number} gunLevel - The level of the shooter.
 * @param {number} range - The range to target.
 * @param {boolean} auto - single or automatic fire mode
 * @param {string} shooterStance - shooter stance
 * @param {string} targetSize - size/stance of target
 * @param {number} shooterMovement - movement of shooter
 * @param {number} targetMovement - movement of target
 * @returns {number} The hit percentage value.
 *
 * @example
 * const eal = calculateEAL(3, 5, false, "Standing", "Running", 0, 0);
 */
export function calculateEAL(
  gunLevel,
  range,
  auto,
  shooterStance,
  targetSize,
  shooterMovement,
  targetMovement
) {
  const sal = salTable[gunLevel];
  const rangeEal = rangeTable[range];

  const stanceEal = shooterStanceMods[shooterStance];
  const sizeEal = targetSizeMods[targetSize][auto ? 1 : 0];

  const movementEal = getMovementMod(shooterMovement, targetMovement, range);

  return sal + rangeEal + stanceEal + sizeEal + movementEal;
}
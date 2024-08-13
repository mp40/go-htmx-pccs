const pelletSpreadTable = Object.freeze({
    "-12": 1,
    "-10": 2,
    "-6": 3,
    "-4": 4,
    "-2": 6,
    0: 8,
    2: 11,
    4: 14,
    6: 19,
    8: 25,
    10: 34,
    12: 45,
    14: 60,
    16: 79,
    18: 100,
});

function getSpreadFromSalm(salm, diceType) {
  let spread = pelletSpreadTable[salm];
  if (diceType === 1000) {
    spread = spread * 10;
  }

  return spread;
}

const getRandomHitLocation = (initialLocation, spread, upperLimit) => {
  const min = initialLocation - spread;
  const max = initialLocation + spread;

  let location = Math.floor(Math.random() * (max - min + 1)) + min;

  if (location < 0 || location > upperLimit) {
    location = getRandomHitLocation(initialLocation, spread, upperLimit);
  }

  return location;
};

export function generateLocations(salm, diceType, hitCount, impactPoint) {
  const spread = getSpreadFromSalm(salm, diceType);

  let locations = [];

  for (let i = 0; i < hitCount; i++) {
    locations.push(getRandomHitLocation(impactPoint, spread, diceType - 1));
  }

  return locations;
}
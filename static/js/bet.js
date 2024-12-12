const betButton = document.getElementById("bet-button");
const betAmount = document.getElementById("bet-amount");
const cashoutButton = document.getElementById("cashout-button");
const oddsDisplay = document.getElementById("odds-display");
const balanceItem = document.getElementById("balance");
const displayBurst = document.getElementById("disp-burst");
const displayCashedOut = document.getElementById("cashed-out");

let waiting = true;
let stake = 1;
let liveBal = 1;
let balance = 10000;
let activeBet = false;
let isGameRunning = false;

function truncateFloat(num, decimalPlaces) {
  const factor = Math.pow(10, decimalPlaces);
  return Math.floor(num * factor) / factor;
}

function processBet() {
  if (activeBet || isGameRunning) {
    if (activeBet) { console.log("has active bet"); }
    return;
  }

  let amount = Number(betAmount.value);
  console.log("placing bet", amount);

  stake = amount;
  activeBet = true;
  balance -= stake;

  updateBalance();
  toggleBetButtons();
}

function updateOdds(odds) {
  let odd = truncateFloat(odds, 2);
  let dp = 0;

  (async function() {
    if (!isGameRunning) {
      isGameRunning = true;
    }
  })();

  oddsDisplay.innerHTML = `${odd.toFixed(2)}<span>X</span>`;

  if (stake < 10) {
    dp = 2;
  } else if (stake < 100) {
    dp = 1;
  } else {
    dp = 0;
  }
  liveBal = truncateFloat((odd * stake), dp);
  cashoutButton.innerHTML = `${liveBal.toFixed(dp)}`;
}

function cashout(burst = false) {
  if (burst) { }
  else if (waiting) {
    balance += Number(stake);
    updateBalance();
  } else if (isGameRunning) {
    balance += Number(liveBal);
    updateBalance();
    displayCashedOut.classList.remove("hidden");
    displayCashedOut.innerHTML = `Cashed Out`;
  }

  if (activeBet) {
    toggleBetButtons();
    activeBet = false;
  }
}

function callStart() {
  isGameRunning = true;
  waiting = false;
  displayCashedOut.classList.add("hidden")
  displayBurst.classList.add("hidden");
  oddsDisplay.classList.replace("text-red-500", "text-green-500")
}

function callBurst() {
  waiting = true;
  cashout(true);

  displayCashedOut.classList.add("hidden");
  displayBurst.classList.remove("hidden");
  oddsDisplay.classList.replace("text-green-500", "text-red-500")
  isGameRunning = false;
  activeBet = false;
}

function updateBalance() {
  balanceItem.innerHTML = balance;
}

function toggleBetButtons() {
  if (betButton.classList.contains("hidden")) {
    betButton.classList.remove("hidden");
    cashoutButton.classList.add("hidden");
  } else {
    betButton.classList.add("hidden");
    cashoutButton.classList.remove("hidden");
  }
}

export {
  callStart,
  callBurst,
  updateOdds,
  cashout,
  processBet,
}

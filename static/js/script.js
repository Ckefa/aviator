import { cashout, callStart, processBet, updateOdds, callBurst } from "./bet.js";

// Attach functions to the global `window` object
window.processBet = processBet;
window.cashout = cashout;

const socket = new WebSocket("ws://localhost:8080/ws"); // Adjust the URL as needed

// When a message is received
socket.onmessage = (event) => {
  const rawData = event.data;
  const data = JSON.parse(rawData);

  if (data.name === "odd") {
    updateOdds(data.msg);
  } else if (data.name === "startrun") {
    callStart();
  } else if (data.name === "burst") {
    console.log("Bursted");
    callBurst();
  } else {
    console.log(data.msg)
  }
}

// When the connection is closed
socket.onclose = (event) => {
  console.log("WebSocket Disconnected:", event.reason);
};

// When an error occurs
socket.onerror = (error) => {
  console.error("WebSocket Error:", error);
};

// Sending a message to the server
socket.onopen = () => {
  console.log("Connection established. Sending message...");
  socket.send("Hello from client");
};


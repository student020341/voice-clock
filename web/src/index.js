import "./index.scss";
import ClockThing from "./stuff/clock.js";

const dom = {
  clock: document.getElementById("clock"),
  speak: document.getElementById("speak-button"), 
  hoursList: document.getElementById("hours-o"),
  minutesList: document.getElementById("minutes-o"),
  ampmList: document.getElementById("ampm-o"),
  customButton: document.getElementById("custom-button"),
  customControls: document.getElementById("custom-controls"),
  autoButton: document.getElementById("auto-announce"),
  autoText: document.getElementById("announce-text"),
  streamerSelect: document.getElementById("streamer-select")
};

// speak kick off speech thing
new ClockThing(dom);

export default class ClockThing {
  constructor(nodes) {
    // loading new audio files
    this.loading = false;
    // audio is playing
    this.playing = false;
    this.nodes = nodes;
    // show custom controls
    this.customTime = false;
    // enable automatic speaking of time at :00
    this.autoAnnounce = false;

    // time display data
    this.time = {
      hours: "-",
      minutes: "-",
      seconds: "-",
      ampm: "-",
    };

    this.timeOverride = {
      hours: "1",
      minutes: "00",
      ampm: "AM"
    };

    // default voice to crumpet
    this.voice = "murdercrumpet";

    // queue of audio players
    this.queue = [];

    this.initStuff();
  }

  // receive click event, start loading files, play when done
  doThing() {
    if (this.loading || this.playing) {
      return;
    }

    // load clips for current time
    const { speak } = this.nodes;
    const { hours, minutes, ampm } = this.time;
    this.loading = true;
    speak.innerText = "loading";

    // normal time or DINK DONK
    const toQueue =
      minutes === "00"
        ? [
            "Dink_Donk",
            "Its",
            hours === "00" ? "12" : hours, // temp: say 12 oclock for midnight
            "oclock",
            ampm,
          ]
        : [
            "Current_Time",
            hours,
            ...(minutes < 10 ? ["oh", minutes] : [minutes]),
            ampm,
          ];

    this.queue = toQueue.map((v) => {
      let stripped = v.toString();
      // 03 will become "oh.ogg" + "3.ogg", need to strip any leading 0s since audio file name does not have them
      if (!isNaN(stripped)) {
        if (stripped.startsWith("0") && stripped !== "0") {
          stripped = stripped.substr(1);
        } else if (stripped === "0") {
          stripped = "12";
        }
      }

      const file = `/audio/${this.voice}/${stripped}.ogg`;
      const player = new Audio(file);
      const ready = new Promise((resolve) => {
        player.addEventListener("canplaythrough", () => resolve());
      });

      return { file, player, ready };
    });

    Promise.all(this.queue.map((q) => q.ready)).then(() => {
      this.loading = false;
      this.playing = true;
      this.nodes.speak.innerText = "playing";
      this.speakTime();
    });
  }

  // todo: refactor into async loop stack machine
  speakTime() {
    if (this.queue.length === 0) {
      this.playing = false;
      this.nodes.speak.innerText = "speak";
      return;
    }
    const current = this.queue.shift();
    console.log(current);
    const { player } = current;
    player.addEventListener("ended", () => this.speakTime());
    player.play();
  }

  // 10 => 10, 60 => 60, 5 => 05
  _pad(n) {
    let num = Number(n);
    if (isNaN(num)) {
      return n;
    } else {
      return num < 10 ? `0${num}` : num;
    }
  }

  // update class time data, update clock text, use AM/PM
  updateClockTime() {
    const { clock } = this.nodes;
    const now = new Date();

    // shitty custom handling :)
    if (this.customTime) {
      const afternoon = this.timeOverride.ampm === "PM";
      const hoursO = afternoon ? Number(this.timeOverride.hours)+12 : this.timeOverride.hours;
      const minutesO = this.timeOverride.minutes;
      now.setHours(Number(hoursO), Number(minutesO), 0);
    }

    this.time.hours = now.getHours();
    this.time.ampm = "AM";
    if (this.time.hours >= 12) {
      this.time.ampm = "PM";
      if (this.time.hours > 12) this.time.hours -= 12;
    }
    this.time.minutes = now.getMinutes();
    if (this.time.minutes < 10) {
      this.time.minutes = `0${this.time.minutes}`;
    }
    this.time.seconds = now.getSeconds();
    if (this.time.seconds < 10) {
      this.time.seconds = `0${this.time.seconds}`;
    }

    // more overrides
    if (this.customTime) {
      this.time.ampm = this.timeOverride.ampm;
    }

    clock.innerText = `${this.time.hours}:${this.time.minutes}:${this.time.seconds} ${this.time.ampm}`;

    // auto announce
    if (this.autoAnnounce && !this.customTime) {
      if (now.getMinutes() === 0 && now.getSeconds() === 0) {
        this.doThing();
        return;
      } else {
        let minutesUntil = this._pad(59 - now.getMinutes());
        let secondsUntil = this._pad(60 - now.getSeconds());
        this.nodes.autoText.innerText = `${minutesUntil} minutes ${secondsUntil} seconds until... 👀`;
      }
    }
    
    if (this.customTime) {
      this.nodes.autoText.innerText = "note - auto announce disabled if custom time is open";
    }
  }

  initStuff() {
    const { hoursList, minutesList } = this.nodes;

    Array.from(new Array(12)).forEach((_, index) => {
      const op = document.createElement("option");
      op.innerText = index + 1;
      hoursList.appendChild(op);
    });

    Array.from(new Array(60)).forEach((_, index) => {
      const op = document.createElement("option");
      op.innerText = index < 10 ? `0${index}` : index;
      minutesList.appendChild(op);
    });

    // event bindings
    setInterval(() => this.updateClockTime(), 500);
    this.nodes.speak.addEventListener("click", () => this.doThing());
    this.nodes.hoursList.addEventListener("change", (event) => this.handleTimeOverride(event, "hours"));
    this.nodes.minutesList.addEventListener("change", (event) => this.handleTimeOverride(event, "minutes"));
    this.nodes.ampmList.addEventListener("change", (event) => this.handleTimeOverride(event, "ampm"));
    this.nodes.customButton.addEventListener("click", () => this.toggleCustomControls());
    this.nodes.autoButton.addEventListener("click", () => this.toggleAutoAnnounce());
    this.nodes.streamerSelect.addEventListener("change", (event) => this.onChangeVoice(event));
  }

  // handle hours/minutes override
  handleTimeOverride(/**@type Event*/ event, w) {
    this.timeOverride[w] = event.target.value;
  }

  // handle on-off style :)
  _toggleOnOff(element, on) {
    if (on) {
      element.classList.remove("orange");
      element.classList.add("green");
    } else {
      element.classList.remove("green");
      element.classList.add("orange");
    }
  }

  _toggleShow(element, on) {
    if (on) {
      element.classList.add("show");
    } else {
      element.classList.remove("show");
    }
  }

  // show / hide custom time controls
  toggleCustomControls() {
    const {customButton, customControls} = this.nodes;
    this.customTime = !this.customTime;
    this._toggleOnOff(customButton, this.customTime);
    this._toggleShow(customControls, this.customTime);
    if (this.customTime) {
      customButton.innerText = "custom time";
    } else {
      customButton.innerText = "custom time";
    }
  }

  // enable auto announce
  toggleAutoAnnounce() {
    const {autoButton, autoText} = this.nodes;
    this.autoAnnounce = !this.autoAnnounce;
    this._toggleOnOff(autoButton, this.autoAnnounce);
    this._toggleShow(autoText, this.autoAnnounce);
  }

  // switch streamer voice - this probably isn't going to scale well
  onChangeVoice(event) {
    this.voice = event.target.value;
  }
}

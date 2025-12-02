// ignoreQueue is used to keep optimistic rendering from duplicating clicks
// I am choosing to do this over server-side tracking of command ownership because it's simpler, and the client isn't authoritative anyway. This is only for rendering purposes
let ignoreQueue = 0;

let ws = new WebSocket(`ws://localhost:8082/pet/${pet_id}/ws`)

if (ws) {
  ws.onopen = () => {
	  console.log("Connected to websocket")
  };

  ws.onmessage = (e) => {
	  console.log(e.data)
	  const event = JSON.parse(e.data)
	  console.table(event)
	  const eventType = event.type
	  switch (eventType) {
		  case "pet":
			  if (ignoreQueue > 0) {
				  ignoreQueue--;
				  return;
			  }
			  increment()
			  break;
		  case "petcount":
			  setCount(Number(event.data))
			  break;
		  case "chat":
			  const msgData = event.data
			  appendChatMessage(msgData.author, msgData.msg)
			  break;
		  default:
			  console.log("Unknown message type")
	  }
  };

  ws.onerror = (e) => {
    // TODO: handle error
  };
  ws.onclose = () => {
    // TODO: handle close
  };
}

// Counter logic (local fallback)
const countEl = document.getElementById('pet-count');
const personalCountEl = document.getElementById('user-pet-count');
const petBtn = document.getElementById('pet-container');

let count = 1000;
function getCount() {
  return count;
}
function setCount(n) {
	count = n;
	updateUI();
}

function updateUI() {
	countEl.innerText = getCount().toLocaleString();
	personalCountEl.innerText = getCount().toLocaleString();
}

function increment() {
  // TODO: implement local increment logic and optionally send to server
	setCount(getCount() + 1);
}

// Pet button logic

// petCommand is always the same, so we declare it once. no memory wasting here.
const petCommand = {
	"type": "pet",
	"data": null
}

function renderOptimistic() {
	increment();
	ignoreQueue++;
}

function pet() {
	countEl.classList.add('count-bump');
	setTimeout(() => countEl.classList.remove('count-bump'), 320);
	renderOptimistic();
	ws.send(JSON.stringify(petCommand))
}

if (petBtn) {
  petBtn.addEventListener('click', pet);
}

// Chat UI
const chatMessages = document.getElementById('chat-message-container');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

function appendChatMessage(user, msg) {
	const msgEl = document.createElement('div');
	msgEl.innerText = `${user}: ${msg}`;

	msgEl.classList.add('message');
	msgEl.classList.add((user === 'You' ? 'me' : 'them'));

	chatMessages.appendChild(msgEl);
}

function sendChat() {
	const msg = chatInput.value;
	// some crisp optimistic rendering
	appendChatMessage('You', msg);

	const msgData = {
		"type": "chat",
		"data": {
			"msg": msg,
			"author": "test"
		}
	}
	ws.send(JSON.stringify(msgData));
	chatInput.value = '';
}

if (chatSend) chatSend.addEventListener('click', sendChat);
if (chatInput) {
  chatInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      sendChat();
    }
  });
}



function setGradientPosition(dog) {

	if (dog == undefined) {
		dog = "daisy"
	}

	const daisyRect = petBtn.getBoundingClientRect();

	const centerX = daisyRect.left + daisyRect.width / 2;
	const centerY = daisyRect.top + daisyRect.height / 2;

	document.body.style.background = `
        radial-gradient(circle at ${centerX}px ${centerY}px, var(--${dog}-gradient-start) 1%, var(--${dog}-gradient-end) 100%
    `;
}

setGradientPosition();

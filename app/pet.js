const ws = new WebSocket(`ws://localhost:8082/pet/${pet_id}/ws`);

ws.onopen = () => {
	console.log("Connected successfully");
};

ws.onmessage = (data) => {
	console.log(data);
};

if (localStorage.getItem("token")) {
	alert("You are already logged in");
	window.location.href = "/";
}


const loginButton = document.getElementById("login-button");
const emailInput = document.getElementById("email-input");
const passwordInput = document.getElementById("password-input");

loginButton.addEventListener("click", async (e) => {
	e.preventDefault();
	await login();
});

async function login() {
	const email = emailInput.value;
	const password = passwordInput.value;

	const res = await fetch("api/login", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ email, password }),
	});

	if (res.ok) {
		data = await res.json();
		console.log(data);

		// Session is now cookie-based; no token storage needed.
		alert("You have been logged in!");
		window.location.href = "/";
	} else {
		console.error("Login failed", await res.text());
	}
}

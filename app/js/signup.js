const emailInput = document.getElementById("email-input");
const passwordInput = document.getElementById("password-input");
const signupButton = document.getElementById("signup-button");

signupButton.addEventListener("click", async (e) => {
	e.preventDefault();
	await signup();
});

async function signup() {
	const email = emailInput.value;
	const password = passwordInput.value;

	const res = await fetch("api/signup", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ email, password }),
	});

	if (res.ok) {
		data = await res.json();
		console.log(data);
		// Session cookie is set by the server; redirect to home
		window.location.href = "/";
	} else {
		console.error("Signup failed", await res.text());
	}
}

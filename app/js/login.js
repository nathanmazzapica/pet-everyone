async function login() {
	const email = emailInput.value;
	const password = passwordInput.value;

	try {
		const res = await fetch("/api/login", {
			method: "POST",
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ email, password }),
		});

		if (!res.ok) {
			const contentType = res.headers.get("content-type") || "";
			if (contentType.includes("application/json")) {
				const errJson = await res.json();
				const msg = errJson?.error || errJson?.message || JSON.stringify(errJson);
				console.error("Login failed", msg);
				alert(`Login failed: ${msg}`);
				return;
			}
			const errorText = await res.text();
			console.error("Login failed", errorText);
			alert(`Login failed: ${errorText || res.status}`);
			return;
		}

		const contentType = res.headers.get("content-type") || "";
		if (contentType.includes("application/json")) {
			const data = await res.json();
			console.log(data);
		}

		// Session is now cookie-based; no token storage needed.
		alert("You have been logged in!");
		window.location.href = "/";
	} catch (err) {
		console.error("Login failed", err);
		alert(`Login failed: ${err?.message || err}`);
	}
}

const loginButton = document.getElementById("login-button");
const emailInput = document.getElementById("email-input");
const passwordInput = document.getElementById("password-input");

loginButton.addEventListener("click", async (e) => {
	e.preventDefault();
	await login();
});

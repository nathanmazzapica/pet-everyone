const displayNameInput = document.getElementById("display-name-input");
const emailInput = document.getElementById("email-input");
const passwordInput = document.getElementById("password-input");
const signupButton = document.getElementById("signup-button");

signupButton.addEventListener("click", async (e) => {
	e.preventDefault();
	await signup();
});

async function signup() {
	const displayName = displayNameInput.value;
	const email = emailInput.value;
	const password = passwordInput.value;

	try {
		const res = await fetch("/api/signup", {
			method: "POST",
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ email, password, display_name: displayName }),
		});

		if (!res.ok) {
			const contentType = res.headers.get("content-type") || "";
			if (contentType.includes("application/json")) {
				const errJson = await res.json();
				const msg = errJson?.error || errJson?.message || res.statusText || JSON.stringify(errJson);
				console.error("Signup failed", msg);
				alert(`Signup failed: ${msg}`);
				return;
			}
			const errorText = await res.text();
			console.error("Signup failed", errorText);
			alert(`Signup failed: ${errorText || res.statusText || res.status}`);
			return;
		}

		const contentType = res.headers.get("content-type") || "";
		if (contentType.includes("application/json")) {
			const data = await res.json();
			console.log(data);
		}
		// Session cookie is set by the server; redirect to home
		window.location.href = "/";
	} catch (err) {
		console.error("Signup failed", err);
		alert(`Signup failed: ${err?.message || err}`);
	}
}

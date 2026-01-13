document.addEventListener("DOMContentLoaded", (event) => {
	if (!localStorage.getItem("token")) {
		alert("You must be logged in to create a pet");
		window.location.href = "/login";
	}
	displayImagePreview();
});

const createForm = document.getElementById('new-pet-form')
const imageInput = document.getElementById('pet-image-input')
const processingPopup = document.getElementById('processing-popup')

imageInput.addEventListener('change', displayImagePreview)

function displayImagePreview() {
	if (!imageInput.files || !imageInput.files[0]) {
		return;
	}
	const file = imageInput.files[0];
	const reader = new FileReader();
	reader.onload = () => {
		const output = document.getElementById('pet-image-preview');
		output.src = reader.result;
	}
	reader.readAsDataURL(file);
}


createForm.addEventListener('submit', e => {
    e.preventDefault()
    uploadPet()
})

function addPeriod() {
	document.getElementById('processing-text').innerText += '.';
}

async function uploadPet() {
    const petName = document.getElementById('pet-name-input').value;
    const petImageFile = document.getElementById('pet-image-input').files[0];

    if (!petImageFile) {
        return;
    }

    const formData = new FormData();
    formData.append('petName', petName);
    formData.append('petImageFile', petImageFile);

	const loadingEffect = setInterval(addPeriod, 500)
    try {
		processingPopup.classList.remove('hidden')
        const res = await fetch(`/api/create`, {
            method: 'POST',
			headers: {
				Authorization: `Bearer ${localStorage.getItem('token')}`,
			},
            body: formData,
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to create pet. Error: ${data.error}`);
        }
        const data = await res.json();
        console.table(data)
        console.log('Successfully created');
        window.location.href = `/pet/${data}`
    } catch(error) {
        alert(`Failed to create pet. Error: ${error}`);
		processingPopup.classList.add('hidden')
		clearInterval(loadingEffect)
		document.getElementById('processing-text').innerText = 'Processing';
    }
}

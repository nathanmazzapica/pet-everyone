document.addEventListener("DOMContentLoaded", () => {
    displayImagePreview();
});

const createForm = document.getElementById('new-pet-form')
const imageInput = document.getElementById('pet-image-input')
const processingPopup = document.getElementById('processing-popup')

if (imageInput) {
    imageInput.addEventListener('change', displayImagePreview)
}

function displayImagePreview() {
    if (!imageInput || !imageInput.files || !imageInput.files[0]) {
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


if (createForm) {
    createForm.addEventListener('submit', e => {
        e.preventDefault()
        uploadPet()
    })
}

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
            credentials: 'include',
            body: formData,
        });
        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Failed to create pet. Error: ${data.error || res.status}`);
        }
        console.table(data)
        console.log('Successfully created');
        window.location.href = `/pet/${data}`
    } catch (error) {
        alert(`Failed to create pet. Error: ${error}`);
    } finally {
        clearInterval(loadingEffect)
        processingPopup.classList.add('hidden')
        document.getElementById('processing-text').innerText = 'Processing';
    }
}

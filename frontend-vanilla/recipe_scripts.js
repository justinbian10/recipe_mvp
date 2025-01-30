const generatePreview = async () => {
	const container = document.querySelector(".recipe-preview-layout");
	url = "http://localhost:8080/v1/recipes/"
	for (let i = 1; i < 5; i++) {
		try {
			const response = await fetch(url + i);
			responseBody = await response.json();
			const recipeElement = createRecipeElement(responseBody.recipe);
			console.log(container);
			console.log(recipeElement);
			container.appendChild(recipeElement);
		} catch(e) {
			console.log(e);
		}
	}
}

const createRecipeElement = (recipe) => {
	const title = recipe.title;
	const imageURL = recipe.image_url;
	const id = recipe.id;
		
	const recipeElement = document.createElement("a");
	recipeElement.setAttribute("href", "/recipe?id=" + id);
	recipeElement.innerHTML = `
		<div class="recipe">
			<img src=${imageURL} />
			<h2>${title}</h2>
		</div>
	`;
	return recipeElement;
}

console.log('hi');
generatePreview();



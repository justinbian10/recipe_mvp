import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';

function RecipeEdit() {
	recipe, setRecipe = useState({});
	searchParams, setSearchParams = useSearchParams();

	useEffect(() => {
		async function startFetching() {
			const apiClient = new ApiClient("http://localhost:8080/v1");
			const recipeObj = await apiClient.getRecipe(searchParams.get("id"));
			setRecipe(recipesObj.recipes);
		}
		startFetching();
	}, [])

	return (
		<>
			{recipe.title}
		</>
	)
}

export default RecipeEdit

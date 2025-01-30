class ApiClient {
	constructor(baseUrl, headers={}) {
		this.baseUrl = baseUrl;
		this.headers = headers;
	}

	async getAllRecipes(filters={}) {
		const url = new URL(`${this.baseUrl}/recipes`);
		console.log(url);
		const response = await fetch(url, {method: 'GET'});
		return this._handleResponse(response);
	}

	async getRecipe(id) {
		const url = new URL(`${this.baseUrl}/recipes/${id}`)
		const response = await fetch(url, {method: 'GET'});
		return this._handleResponse(response);
	}

	async getAllIngredients(name) {
		const url = new URL(`${this.baseUrl}/ingredients?name=${name}`)
		const response = await fetch(url, {method: 'GET'});
		return this._handleResponse(response);
	}

	async _handleResponse(response) {
		if (!response.ok) {
			const errorData = await response.json();
			throw new Error(errorData.message || 'Something went wrong');
		}
		return response.json()
	}
}

export default ApiClient

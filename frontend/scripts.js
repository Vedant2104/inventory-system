let productlist = document.getElementById("products");
let button = document.getElementById("button");
function addProducts(product){
    let products = document.getElementById("product-list");
    let productDiv = document.createElement("div");
    let name = document.createElement("h3");
    name.textContent = product.name
    let desc = document.createElement("p");
    desc.textContent = product.description
    let price = document.createElement("p");
    price.textContent = `Price: $${product.price}` 
    productDiv.appendChild(name);
    productDiv.appendChild(desc);
    productDiv.appendChild(price);
    productDiv.setAttribute("class", "product");
    products.appendChild(productDiv);
}

let fetched = false

function fetchProducts(){
    if (fetched){
        return
    }
    fetched = true
    console.log("hello")
    fetch("http://localhost:8080/product")
    .then(res => res.json())
    .then(data => {
        data.forEach(product => addProducts(product))
    })
    .catch(err => console.log(err))
}

button.addEventListener("click", fetchProducts)

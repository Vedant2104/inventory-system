import ProductCard from "./ProductCard";
import "../styles/Product_list.css";
import { useEffect, useState } from "react";
import type { Product, Category } from "../interfaces/inventory";
import Loader from "./Loader";
import { useSearchParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";

function ProductList() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  // const [categoryFilter, setCategoryFilter] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error , setError] = useState<string>("")
  const [searchParams , setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  const categoryFilterParam = searchParams.get("category") || "";
  function HandleFetchProducts() {
    setIsLoading(true);

    fetch(`http://localhost:8080/product?category=${categoryFilterParam}`, {
      method: "GET",
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error("Failed to fetch products");
      })
      .then((data) => {
        setProducts(data || []);
      })
      .catch((err) => {
        console.log(err);
        setError(err.message)
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  function HandleFetchCategories() {
    setIsLoading(true);

    fetch("http://localhost:8080/category", {
      method: "GET",
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error("Failed to fetch categories");
      })
      .then((data) => {
        setCategories(data || []);
        console.log(data);
      })
      .catch((err) => {
        console.log(err);
        setError(err.message)
      })
      .finally(() => {
        setIsLoading(false);
      });
  }
  useEffect(() => {
    HandleFetchCategories();
  }, []);


  useEffect(() => {
    HandleFetchProducts();
  }, [categoryFilterParam]);

  if (isLoading) {
    return <Loader/>
  }
  if (error) {
    return <p className="error">{error}</p>
  }

  return (
    <div className="product-list">
      
      <button className="add-item" onClick={() => navigate("/product/new")}>➕Add Product</button>

      <div className="category-filter-container">
        <label>Filter by category:</label>
        <select
          className="category-filter"
          value={categoryFilterParam}
          onChange={(e) => {
            // setCategoryFilter(e.target.value);
            // setSearchParams({ category: e.target.value });
            setSearchParams((prev)=>{
              const newParam = new URLSearchParams(prev);
              newParam.set("category", e.target.value);
              return newParam
            })
          }}
        >
          <option key="all" value="">
            All
          </option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      {products.length === 0 ? (
        <p>No products found for this category</p>
      ) : (
        products.map((item) => <ProductCard key={item.id} {...item} />)
      )}
    </div>
  );
}

export default ProductList;

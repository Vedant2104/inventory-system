import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import type { Product, Category } from "../interfaces/inventory";
import "../styles/ProductPage.css";

export default function ProductPage() {
  const { id } = useParams<{ id: string }>();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [product, setProduct] = useState<Product | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [error, setError] = useState<string>("");
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [formData, setFormData] = useState<Product | null>(null);

  function handleFetchProduct() {
    setIsLoading(true);
    fetch(`http://localhost:8080/product/${id}`, { method: "GET" })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch product");
        return res.json();
      })
      .then((data) => {
        setProduct(data);
        setFormData(data);
      })
      .catch((err) => setError(err.message))
      .finally(() => setIsLoading(false));
  }

  function handleFetchCategories() {
    fetch("http://localhost:8080/category", { method: "GET" })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch categories");
        return res.json();
      })
      .then((data) => setCategories(data || []))
      .catch((err) => setError(err.message));
  }

  function handleEditClick() {
    handleFetchCategories();
    setIsEditing(true);
  }

  function handleCancel() {
    setFormData(product);
    setIsEditing(false);
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) {
    const { name, value } = e.target;

    if (name === "category") {
      const selected = categories.find((c) => c.id === value);
      setFormData((prev) => prev ? { ...prev, category: selected! } : prev);
    }else if(name == 'price' || name == 'quantity'){
      setFormData((prev) => prev ? { ...prev, [name]: Number(value) } : prev);
    }
    else {
      setFormData((prev) => prev ? { ...prev, [name]: value } : prev);
    }
  }

  function handleSave() {
    setIsLoading(true);
    const payload = {
        name : formData?.name,
        description : formData?.description,
        category : formData?.category.id,
        price : formData?.price,
        brand : formData?.brand,
        quantity : formData?.quantity
    }
    fetch(`http://localhost:8080/product/${id}`,{
        method: "PATCH",
        headers:{
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    }).then(res=>{
        if(!res.ok) throw new Error("Failed to update product");
        return res.json()
    }).then(data=>{
        setProduct(data);
    }).catch(err=>setError(err.message))
    .finally(()=>{
        setIsEditing(false);
        setIsLoading(false);
    })
  }

  useEffect(() => {
    handleFetchProduct();
  }, []);

  if (isLoading) return <div className="page-state">Loading...</div>;
  if (error)     return <div className="page-state error">{error}</div>;
  if (!product)  return <div className="page-state">No product found.</div>;

  return (
    <div className="product-page">
      <div className="product-card">

        {/* ---- VIEW MODE ---- */}
        {!isEditing && (
          <>
            <div className="card-header">
              <h1 className="product-name">{product.name}</h1>
              <Link to = {`/?category=${product.category.id}`}><span className="product-category">{product.category.name}</span></Link>
            </div>

            <p className="product-description">{product.description}</p>

            <div className="product-meta">
              <div className="meta-item">
                <span className="meta-label">Price</span>
                <span className="meta-value">${product.price}</span>
              </div>
              <div className="meta-item">
                <span className="meta-label">Brand</span>
                <span className="meta-value">{product.brand}</span>
              </div>
              <div className="meta-item">
                <span className="meta-label">Quantity</span>
                <span className="meta-value">{product.quantity}</span>
              </div>
            </div>

            <button className="btn-edit" onClick={handleEditClick}>
              Edit Product
            </button>
          </>
        )}

        {/* ---- EDIT MODE ---- */}
        {isEditing && formData && (
          <>
            <div className="card-header">
              <h1 className="product-name">Edit Product</h1>
            </div>

            <div className="form-group">
              <label>Name</label>
              <input name="name" value={formData.name} onChange={handleChange} />
            </div>

            <div className="form-group">
              <label>Description</label>
              <textarea name="description" value={formData.description} onChange={handleChange} />
            </div>

            <div className="form-group">
              <label>Category</label>
              <select name="category" value={formData.category.id} onChange={handleChange}>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Price</label>
                <input name="price" type="number" value={formData.price} onChange={handleChange} />
              </div>
              <div className="form-group">
                <label>Brand</label>
                <input name="brand" value={formData.brand} onChange={handleChange} />
              </div>
              <div className="form-group">
                <label>Quantity</label>
                <input name="quantity" type="number" value={formData.quantity} onChange={handleChange} />
              </div>
            </div>

            <div className="form-actions">
              <button className="btn-cancel" onClick={handleCancel}>Cancel</button>
              <button className="btn-save" onClick={handleSave}>Save</button>
            </div>
          </>
        )}

      </div>
    </div>
  );
}
import { useParams } from "react-router-dom"
import type { Category } from "../interfaces/inventory";
import { useState, useEffect } from "react";
import Loader from "./Loader";
import "../styles/CategoryPage.css"

export default function CategoryPage() {
    const { id } = useParams<{ id: string }>();
    const [category, setCategory] = useState<Category | null>(null);
    const [formData, setFormData] = useState<Category | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>("");
    const [isEditing, setIsEditing] = useState<boolean>(false);

    function handleFetchCategory() {
        setIsLoading(true);
        fetch(`http://localhost:8080/category/${id}`, { method: "GET" })
            .then(resp => {
                if (!resp.ok) throw new Error("Failed to fetch category");
                return resp.json();
            })
            .then(data => {
                setCategory(data);
                setFormData(data);
            })
            .catch(err => setError(err.message))
            .finally(() => setIsLoading(false));
    }

    function handleSave(e: React.SyntheticEvent<HTMLFormElement>) {
        e.preventDefault();
        setIsLoading(true);
        fetch(`http://localhost:8080/category/${id}`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(formData)
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to update category");
                return res.json();
            })
            .then(data => {
                setCategory(data);
                setFormData(data);
            })
            .catch(err => setError(err.message))
            .finally(() => {
                setIsEditing(false);
                setIsLoading(false);
            });
    }

    function handleCancel() {
        setFormData(category);
        setIsEditing(false);
    }

    useEffect(() => {
        handleFetchCategory();
    }, []);

    if (isLoading) return <Loader />;
    if (error) return <p className="page-state error">{error}</p>;
    if (!category) return <p className="page-state">Category not found</p>;

    return (
        <div className="category-page">
            <div className="category-card">

                {!isEditing && (
                    <>
                        {/* top accent banner with first letter */}
                        <div className="category-banner">
                            <span className="category-initial">{category.name.charAt(0)}</span>
                        </div>

                        <div className="category-body">
                            <div className="category-header">
                                <div>
                                    <p className="category-label-small">Category</p>
                                    <h1 className="category-name">{category.name}</h1>
                                </div>
                                <span className="category-id-badge">ID #{id}</span>
                            </div>

                            <div className="category-section">
                                <p className="section-label">Description</p>
                                <p className="category-description">
                                    {category.description || "No description provided."}
                                </p>
                            </div>

                            <button className="btn-edit" onClick={() => setIsEditing(true)}>
                                Edit Category
                            </button>
                        </div>
                    </>
                )}

                {isEditing && formData && (
                    <div className="category-body">
                        <div className="category-header">
                            <h1 className="category-name">Edit Category</h1>
                        </div>

                        <form className="category-form" onSubmit={handleSave}>
                            <div className="form-group">
                                <label>Name</label>
                                <input
                                    type="text"
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                />
                            </div>

                            <div className="form-group">
                                <label>Description</label>
                                <textarea
                                    value={formData.description}
                                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                />
                            </div>

                            <div className="form-actions">
                                <button type="button" className="btn-cancel" onClick={handleCancel}>Cancel</button>
                                <button type="submit" className="btn-save">Save</button>
                            </div>
                        </form>
                    </div>
                )}

            </div>
        </div>
    );
}
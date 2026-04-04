import {  useParams } from "react-router-dom"
import type { Category } from "../interfaces/inventory";
import { useState , useEffect} from "react";
import Loader from "./Loader";
import "../styles/CategoryPage.css"
export default function CategoryPage() {
    const {id} = useParams<{id: string}>();
    const [category ,setCategory] = useState<Category | null>(null);
    const [formData ,setFormData] = useState<Category | null>(null);
    const [isLoading , setIsLoading] = useState<boolean>(false);
    const [error , setError] = useState<string>("");
    const [isEditing , setIsEditing] = useState<boolean>(false);
    function handleFetchCategory(){
        setIsLoading(true);

        fetch(`http://localhost:8080/category/${id}`, { method: "GET" })
        .then(resp=>{
            if(!resp.ok){
                throw new Error("Failed to fetch category")
            }
            return resp.json();
        })
        .then(data=>{
            setCategory(data);
            setFormData(data);
        })
        .catch(err=>{
            console.log(err);
            setError(err.message);
        })
        .finally(()=>{
            setIsLoading(false);
        })
    }

    function handleSave(e: React.SyntheticEvent<HTMLFormElement>){
        e.preventDefault(); 
        setIsLoading(true);
        fetch(`http://localhost:8080/category/${id}`,{
            method: "PATCH",
            headers:{
                "Content-Type": "application/json"
            },
            body: JSON.stringify(formData)
        }).then(res=>{
            if(!res.ok) throw new Error("Failed to update category");
            return res.json()
        }).then(data=>{
            setCategory(data);
            setFormData(data);
        }).catch(err=>setError(err.message))
        .finally(()=>{
            setIsEditing(false);
            setIsLoading(false);
        })
    }

    function handleCancel(){
        setFormData(category);
        setIsEditing(false);
    }

    useEffect(()=>{
        handleFetchCategory();
    },[])

    if(isLoading){
        return <Loader />;
    }
    if(error){
        return <p>{error}</p>
    }
    if(!category){
        return <p>Category not found</p>
    }

  return (
    <div className="category-page">
        <div className="category-card">

            {!isEditing && (
                <>
                    <div className="category-header">
                        <h1 className="category-name">{category.name}</h1>
                    </div>
                    <hr/>
                    <p className="category-description">
                        {category.description}
                    </p>

                    <button className="edit-category" onClick={()=>{setIsEditing(true)}}>Edit Category</button>
                </>
            )}

            {isEditing && formData && (
                <form className="category-form" onSubmit={(e)=>{handleSave(e)}}>
                    <label className="category-label">Category Name:</label>   
                    <input className="category-input" type="text" value={formData.name} onChange={(e)=>{setFormData({...formData , name: e.target.value})}} />

                    <label className="category-label">Category Description:</label>
                    <textarea className="category-input" value={formData.description} onChange={(e)=>{setFormData({...formData , description: e.target.value})}} />

                    {/* <button className="save-category" onClick={handleSave}>Save</button> */}
                    <button type="submit" className="save-category">Save</button>
                    <button className="cancel-edit" onClick={()=>{handleCancel()}}>Cancel</button>
                </form>
            )}
        </div>
    </div>
  )
}

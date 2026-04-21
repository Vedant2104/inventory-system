import { useState } from "react"
import type { ProductCountByCategory } from "../../interfaces/inventory"
import Loader from "../Loader";
import { Link } from "react-router-dom";
import "../../styles/CountByCategoryReport.css"
import { CSVDownloader } from "../../utils/CSVDowloader";
export default function ProductCountByCategoryReport() {
    const [productCountByCategory, setProductCountByCategory] = useState<ProductCountByCategory[]>([])
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>("");
    const [minValue, setMinValue] = useState<string>("0");
    const [maxValue, setMaxValue] = useState<string>("10");
    const [hasGenerated, setHasGenerated] = useState<boolean>(false);

    function handleFetch() {
        if (Number(minValue) > Number(maxValue)) {  // guard: min can't exceed max
            setError("Min value cannot be greater than max value");
            return;
        }

        setIsLoading(true);
        setHasGenerated(false);
        setProductCountByCategory([]);
        setError("");

        fetch(`http://localhost:8080/product/report/countbycategory?minValue=${minValue}&maxValue=${maxValue}`, {
            method: "GET"
        })
        .then(async res => {
            if (!res.ok) {
                const errText = await res.text();
                throw new Error(errText || "Failed to fetch report");
            }
            return res.json();
        })
        .then(data => {
            setProductCountByCategory(data || []);
        })
        .catch(err => {
            setError(err.message);
        })
        .finally(() => {
            setIsLoading(false);
            setHasGenerated(true);
        });
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") handleFetch();
    }

    if (isLoading) return <Loader />;

    return (
        <div className="report-page">
            <button className={`btn-download ${!hasGenerated ? "disabled" : ""}`} disabled={!hasGenerated} onClick={() => CSVDownloader(productCountByCategory , "product-count-by-category")}>Download</button>


            <div className="report-container">

                <div className="report-header">
                    <h2 className="report-title">Product Count by Category</h2>
                    <p className="report-subtitle">Number of products per category within a count range</p>
                </div>

                <div className="report-controls">
                    <div className="input-group">
                        <label htmlFor="minValue">Min Count</label>
                        <input
                            id="minValue"
                            type="number"
                            value={minValue}
                            min="0"
                            onChange={(e) => setMinValue(e.target.value)}
                            onKeyDown={handleKeyDown}
                        />
                    </div>
                    <div className="input-group">
                        <label htmlFor="maxValue">Max Count</label>
                        <input
                            id="maxValue"
                            type="number"
                            value={maxValue}
                            min="0"
                            onChange={(e) => setMaxValue(e.target.value)}
                            onKeyDown={handleKeyDown}
                        />
                    </div>
                    <button className="btn-generate" onClick={handleFetch}>
                        Generate Report
                    </button>
                </div>

                {error && <p className="report-error">{error}</p>}

                {hasGenerated && productCountByCategory.length === 0 && !error && (
                    <p className="report-empty">No categories found in this range.</p>
                )}

                {productCountByCategory.length > 0 && (
                    <>
                        <p className="report-count">{productCountByCategory.length} category(s) found</p>
                        <table className="report-table">
                            <thead>
                                <tr>
                                    <th>Category</th>
                                    <th>Product Count</th>
                                </tr>
                            </thead>
                            <tbody>
                                {productCountByCategory.map((item) => (
                                    <tr key={item.category_id}>
                                        <td>
                                            <Link to={`/category/${item.category_id}`}>{item.category}</Link>
                                        </td>
                                        <td>
                                            <span className="count-badge">{item.count}</span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </>
                )}

            </div>
        </div>
    );
}
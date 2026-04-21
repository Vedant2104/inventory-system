import type { PriceSegmentation } from "../../interfaces/inventory";
import { useEffect, useState } from "react";
import "../../styles/CountByCategoryReport.css";
import { CSVDownloader } from "../../utils/CSVDowloader";
import Loader from "../Loader";
import { Link } from "react-router-dom";

export default function PriceSegmentationReport() {
  const [priceSegmentation, setPriceSegmentation] = useState<
    PriceSegmentation[]
  >([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  function handleFetch() {
    setLoading(true);
    setError("");
    fetch("http://localhost:8080/product/report/pricesegmentation", {
      method: "GET",
    })
      .then(async (res) => {
        if (!res.ok) {
          const errText = await res.text();
          throw new Error(errText || "Failed to fetch report");
        }
        return res.json();
      })
      .then((data) => {
        setPriceSegmentation(data || []);
      })
      .catch((err) => {
        setError(err.message);
        console.log(err);
      })
      .finally(() => {
        setLoading(false);
      });
  }

  useEffect(() => {
    handleFetch();
  }, []);

  if (loading) {
    return <Loader />;
  }

  return (
    <div className="report-page">
      <button
        className={`btn-download ${
          !priceSegmentation.length ? "disabled" : ""
        }`}
        disabled={!priceSegmentation.length}
        onClick={() =>
          CSVDownloader(priceSegmentation, "PriceSegmentationReport")
        }
      >
        Download
      </button>

      <div className="report-container">
        <div className="report-header">
          <h2 className="report-title">Product Price Segmentation</h2>
          <p className="report-subtitle">
            Number of products per category within a price range
          </p>
        </div>

        {error && <p className="report-error">{error}</p>}

        {!priceSegmentation.length && (
          <p className="report-empty">No data to display</p>
        )}

        {priceSegmentation.length > 0 && (
          <>
            <p className="report-count">
              {priceSegmentation.length} category(s) found
            </p>
            <table className="report-table">
              <thead>
                <tr>
                  <th>Category</th>
                  <th>Budget Range (less than $2000)</th>
                  <th>Middle Range (Between $2000 and $5000)</th>
                  <th>Premium Range (more than $5000)</th>
                </tr>
              </thead>
              <tbody>
                {priceSegmentation.map((item) => (
                  <tr key={item.category_id}>
                    <td>
                      <Link to={`/category/${item.category_id}`}>
                        {item.category}
                      </Link>
                    </td>
                    <td>
                      <span className="count-badge">{item.budget}</span>
                    </td>
                    <td>
                      <span className="count-badge">{item.midRange}</span>
                    </td>
                    <td>
                      <span className="count-badge">{item.premium}</span>
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

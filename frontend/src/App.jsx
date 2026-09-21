import { useEffect, useState } from "react";

function App() {
  const API_URL = import.meta.env.VITE_API_URL;
  const [customers, setCustomers] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch(`${API_URL}/api/customers`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("顧客情報の取得に失敗しました");
        }

        return response.json();
      })
      .then((data) => {
        setCustomers(data);
      })
      .catch((error) => {
        console.error(error);
        setError(error.message);
      });
  }, []);

  return (
    <div>
      <h1>KT CRM Portfolio</h1>

      <h2>顧客一覧</h2>

      {error && <p>{error}</p>}

      <ul>
        {customers.map((customer) => (
          <li key={customer.id}>
            {customer.name} / {customer.email} / {customer.company}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
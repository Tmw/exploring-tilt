import { useState } from "react";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <div>
        <strong>Count is:{count}</strong>
      </div>
      <div>
        <button onClick={() => setCount(count + 1)}>Inc</button>
        <button onClick={() => setCount(count - 1)}>Dec</button>
      </div>
    </>
  );
}

export default App;

import React, { useState, useEffect } from 'react';
import { Link, Stack, Grid, Row, Column } from "react-ui";
import './App.css';
import SecuritiesTable from './Graphs';
import axios from 'axios';

const client = axios.create({
  baseURL: "/securities" 
});

function App() {
  const [securities, setSecurities] = useState([]);

  // GET with Axios
  useEffect(() => {
    const fetchSecurities = async () => {
      let response = await client.get();

      setSecurities(response.data.securities);
    };
    fetchSecurities();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Market Analysis</h1>
      </header>
      <Grid>
        <Row />
        <Column span={2}>
          <Stack direction="vertical">
            <Link
              size={20}
              href=""
              target="_blank"
            >Securities</Link>
            <Link size={20}>Graphs</Link>
          </Stack>
        </Column>
        <Column span={10}>
          <h3>Securities</h3>
          <SecuritiesTable securities={securities}/>
        </Column>
      </Grid>
    </div>
  );
}

export default App;

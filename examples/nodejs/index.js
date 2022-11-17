const express = require('express');
const dotenv = require('dotenv');


const app = express()

app.get('/', (req, res) => {
    res.send('Hello World!');
});

app.get('/env', (req, res) => {
    const name = req.query.name;
    if (!name) {
        res.statusCode = 400;
        res.send(`env name required.`);
    } else {
        res.send(process.env[name.toUpperCase()]);
    }
});

app.post('/env/notify', (req, res) => {
    req.on('data', (chunk) => {
        console.log(chunk.toString());
        res.send('ok');
    })
});

let server;

function start() {
    dotenv.config()
    server = app.listen(3000, () => {
        console.log(`Example app listening on port 3000, pid=${process.pid}`);
    });
}

function stop() {
    if (!server) return;

    server.close();
}


process.on('SIGUSR1', () => {
    console.log('Restarting express server...');
    stop();
    start();
});

start();
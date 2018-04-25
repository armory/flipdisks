import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';
import * as _ from 'lodash';
import * as request from 'request-promise';


const config = {
	server: {
		method: 'http',
		hostname: 'localhost',
		port: 8080,
		get baseUrl() {
			return `${config.server.method}://${config.server.hostname}:${config.server.port}`;
		}
	}
};


class App extends Component {
	constructor() {
		super();
	}


	async componentDidMount() {
		let frameNumber = 0;

		setInterval(async () => {
			try {
				const playlist = await request(`${config.server.baseUrl}/v1/sites/armory/playing`, {json: true});
				const video = playlist.videos[0];


				const state = {
					layout: video.layout,
				};

				_.flatten(video.layout).forEach(boardName => {
					state[boardName] = video.frames[frameNumber][boardName];
				});

				this.setState(state);

				frameNumber = (frameNumber + 1) % video.frames.length;
			} catch(e) {

			}
		}, 1000);
	}


	render() {
		if (!this.state) {
			return (<div></div>)
		}

		return (
			<div className="App">
				<header className="App-header">
					<h1 className="App-title">Flip Disc Simulator</h1>
				</header>
				<div className="board-container">

					{this.state.layout.map(yAxisBoardNames => {
						return yAxisBoardNames.map((name, index) => {
							return (<div key={index}>
								{this.state[name].map((row, index) => {
									return (<div key={index}>{row.map(dot => <span class="flipdot">{dot}</span>)}</div>);
								})}</div>);
						})
					})}
				</div>

			</div>
		);
	}
}

export default App;

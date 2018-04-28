import React, {Component} from 'react';
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


export class NowPlaying extends Component {
	constructor() {
		super();
	}


	async componentDidMount() {
		const playlist = await request(`${config.server.baseUrl}/v1/sites/armory/playing`, {json: true});
		const video = playlist.videos[0];

		let frameNumber = 0;

			setInterval(() => {
				const state = {
					layout: video.layout,
					boards: {}	// each board per ferm
				};

				// set each board into the state
				_.flatten(video.layout).forEach(boardName => {
					state.boards[boardName] = video.frames[frameNumber][boardName];
				});

			_.flatten(video.layout).forEach(boardName => {
				state[boardName] = video.frames[frameNumber][boardName];
			});

			this.setState(state);

			frameNumber = (frameNumber + 1) % video.frames.length;
		}, 1000 / video.fps);
	}

	render() {
		const drawABoard = (boardName, columnIndex) => {
			// for each row, create a span of columns
			return (
				<div className="board-row" key={columnIndex}>
					{
						this.state.boards[boardName].map((boardRow, boardIndex) => {
							return (
								<div key={boardIndex}>{boardRow.map(dot => <span className="a-single-flipdisk">{dot}</span>)}</div>);
						})
					}
				</div>);
		};

		if (!this.state) {
			return (<div></div>)
		}

		return (
			<div className="App">
				<header className="App-header">
					<h1 className="App-title">Flip Disc Simulator</h1>
				</header>
				<div className="board-container">

					{this.state.layout.map(rowLayout => {
						// create a row of boards by column
						return rowLayout.map((boardName, columnIndex) => {
							return drawABoard(boardName, columnIndex)
						})
					})}
				</div>

			</div>
		);
	}
}

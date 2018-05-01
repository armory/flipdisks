import React, { Component } from 'react';
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
		const board = await request.post({
			url: `${config.server.baseUrl}/v1/fonts/render`,
			json: true,
			body: {
				"fontName": "TI84",
				"text": "good job",
				"spaceWidth": 4,
				"kerning": 0
			}
		});

		console.log(board)
		this.setState({
			board: board
		});
	}

	render() {
		const drawABoard = () => {
			// for each row, create a span of columns
			return (
				<div className="board-container">
					{
						this.state.board.map((letter, letterIndex) => {
							return (
								<span>
									{
										letter.map((boardRow, boardIndex) => {
											return (
												<div className="board-row" key={boardIndex}>
													{boardRow.map(dot => <span className="a-single-flipdisk">{dot}</span>)}
												</div>
											);
										})
									}
								</span>
							)})
					}
				</div>
			);
		};

		if (!this.state) {
			return (<div></div>)
		}

		return (
			<div className="App">
				<header className="App-header">
					<h1 className="App-title">Flip Disc Simulator</h1>
				</header>
				<div className="">
					{drawABoard()}
				</div>

			</div>
		);
	}
}

export default App;

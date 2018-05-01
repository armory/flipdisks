import React, { Component } from 'react';
import './App.css';
import * as _ from 'lodash';
import * as request from 'request-promise';
import { debug } from 'util';
// import Font from '/Users/jcheng305/go/src/github.com/armory/flipdisks/webclient/src/Font.js';
// import DropDown from '/Users/jcheng305/go/src/github.com/armory/flipdisks/webclient/src/DropDown.js';

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
	constructor(props) {
		super(props);

		this.state = {
			text: "Hello",
			board: null
		}

		this.handleKeyPress = this.handleKeyPress.bind(this);
	}



	async componentDidMount() {
 		const board = await request.post({
			url: `${config.server.baseUrl}/v1/fonts/render`,
			json: true,
			body: {
				"fontName": "TI84",
				"text": "",
				"spaceWidth": 0,
				"kerning": 0
			}
		});

		this.setState({
			// text: "Test",
			board: board
		});
		console.log(board)
	}

	handleKeyPress(event) {
		this.setState({
			text: event.target.value
		});
		
		return request.post({
			url: `${config.server.baseUrl}/v1/fonts/render`,
			json: true,
			body: {
				"fontName": "TI84",
				"text": event.target.value,
				"spaceWidth": 0,
				"kerning": 0
			}
		})
		.then((board) => {
			this.setState({
				board:board
			})
		})
		
		console.log(2)
				// this.setState({
		// 	board: board
		// });
		console.log(3)
	}


	render() {
		const drawABoard = () => {
			if (!this.state.board) {
				return (<div></div>)
			}
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

				<div className="text-box-container">
					Font Type<select></select>
					Grid Size<select></select>
					Font Size<select></select>
				</div>
				<input type="text" onChange={this.handleKeyPress} value={this.state.text}></input>
				<div className="">
					{drawABoard()}
				</div>
				{/* <Font/> */}
				{/* <DropDown/> */}
			</div>
		);
	}
}

export default App;
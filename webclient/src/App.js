import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';
import * as _ from 'lodash';



const payload = require('./payload.json');


class App extends Component {
	constructor() {
		super();
	}


	componentDidMount() {
		let frameNumber = 0;

		setInterval(() => {
			const state = {
				layout: payload.layout,
			};

			_.flatten(payload.layout).forEach(boardName => {
				state[boardName] = payload.frames[frameNumber][boardName];
			});

			this.setState(state);

			frameNumber = (frameNumber + 1) % payload.frames.length;
		}, payload.frameRate);
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
						return yAxisBoardNames.map(name => {
							return (<div>
								{this.state[name].map(x => {
									return (<div>{x.join(' ')}</div>);
								})}</div>);
						})
					})}
				</div>

			</div>
		);
	}
}

export default App;

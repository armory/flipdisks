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
					<img src={logo} className="App-logo" alt="logo"/>
					<h1 className="App-title">Welcome to React</h1>
				</header>
				<pre>
					{this.state.layout.map(yAxisBoardNames => {
						return yAxisBoardNames.map(name => {
							return (<div>{name}
								{this.state[name].map(x => {
									return (<div>{x.join(' ')}</div>);
								})}</div>);
						})
					})}
				</pre>

			</div>
		);
	}
}

export default App;


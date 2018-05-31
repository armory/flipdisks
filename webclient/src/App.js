import React, { Component } from 'react';
import './App.css';
import { debug } from 'util';
import Jumbo from './Jumbotron'
import Font from './Font'



class App extends Component {
	render() {
		return (				
			<div>
				<div className="home-page">
                  <Jumbo/>
			    </div>

				<div className="font-page" id="font">
                    {/*<Font/>*/}
				</div>
				</div>

		);
	}
}

export default App;
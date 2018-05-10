import React, { Component } from 'react';
import './App.css';
import * as _ from 'lodash';
import * as request from 'request-promise';
import { debug } from 'util';
import Font from './Font';
import Typist from 'react-typist';
import {Jumbotron, Button} from 'react-bootstrap';
import GithubCorner from 'react-github-corner';
import Slider from './Slider';

// import MyNavbar from './Navbar';


class App extends Component {
	render() {
		return (				
			<div>

				<GithubCorner  href="https://github.com/armory/flipdisks" />

				<div className="home-page">
				<Jumbotron className="jumbotron1">
					<Typist>
						<h1>Hello, world!</h1>
					<p>
						This is a simple flip disc simulator. =D 
  					</p>
					<p>
						One may ask,what is a flip disc exactly?
					</p>
					<Typist.Delay ms={500} />
					<p>
						The flip-disc display consists of a grid of small metal discs that are black on one side and a bright color on the other (typically white or day-glo yellow), set into a black background. With power applied, the disc flips to show the other side. Once flipped, the discs will remain in position without power.
					</p>
					<p>
						<Button href="#font" bsStyle="primary">Try out our demo!</Button>
					</p>
					</Typist>
				</Jumbotron>
			</div>
				<div className="font-page" id="font">
					<Font/>
				</div>
				<div className="" id="">
				SLIDER SHOULD GO HERE:
				<Slider/>
				</div>
			</div>
		);
	}
}

export default App;
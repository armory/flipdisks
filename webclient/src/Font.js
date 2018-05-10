import React from 'react';
import './App.css';
import * as request from 'request-promise';
import ReactDynamicFont from 'react-dynamic-font';
import Slider from 'carbon-components';
import FontSlider from './Slider';

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
export default class Font extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            text: " ",
            board: null,
            letter: " "
        }

        this.handleKeyPress = this.handleKeyPress.bind(this)
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
            board: board,
            letter: board.letter
        });
        console.log(board)
        console.log(board.letter)
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
                "spaceWidth": 1,
                "kerning": 0
            }
        })
            .then((board) => {
                this.setState({
                    board: board
                })
            })
        console.log(2)
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
                    BOARD CONTAINER(THAT SHOULD RENDER REACT DYNAMIC LETTERS AS A STRING):
                        {/* <ReactDynamicFont smooth content={this.state.board + "&#10;" } /> */}
                        <ReactDynamicFont smooth whiteSpace = {"wrap"} content={this.state.board } />

                    {
                        this.state.board.map((line, lineNumber) => {
                            return (<div className="line-container">
                            LINE CONTAINER:
                                {line.map((letter, letterIndex) => {
                                    return (
                                        <span className="letters">
                                        LETTER SPAN:
                                            {
                                                letter.map((boardRow, boardIndex) => {
                                                    return (
                                                        <div className="board-row" key={boardIndex}>
                                                            BOARD ROW:
                                                            {boardRow.map(dot => <span className="a-single-flipdisk"> SINGLE DOT:{dot}</span>)} 

                                                        </div>
                                                    );
                                                })
                                            }
                                        </span>
                                    )
                                })}
                            </div>)
                        })
                    }

                </div>
            );
        };

        if (!this.state) {
            return (<div></div>)
        }
        return (
            <div className= "font-container">
            <div>
                <h1> Try it out here! Type something in the box below! </h1>
                    SLIDER SHOULD GO HERE:
				    <FontSlider />
            </div>
            <div>
                
            <textarea type="text"
            onChange={this.handleKeyPress} 
            value={this.state.text}>

            </textarea>
            </div>
            Text Container is here:
            <div className="text-container-with-padding">
           <h1> REACT DYNAMIC FONT RESULT HERE: </h1>
            <ReactDynamicFont smooth content = {this.state.text} />
                    {drawABoard()}
            </div>
            </div>
        )
    }
}


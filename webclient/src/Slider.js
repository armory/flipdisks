import React from 'react';
import Slider from 'carbon-components';

export default class FontSlider extends React.Component {

    render() {
        return (
            <div class="bx--form-item">
                <div style={{ marginTop: "2rem" }}>
                    <Slider
                        id="slider"
                        value={50}
                        min={0}
                        max={100}
                        labelText="Slider Label"
                        onChange={onChange()}
                    />
                </div>
                <label for="slider" class="bx--label">Slider Label</label>
                <div class="bx--slider-test">
                    <div class="bx--slider-container">
                        <span class="bx--slider__range-label">0</span>
                        <div class="bx--slider" data-slider data-slider-input-box="#slider-input-box">
                            <div class="bx--slider__track"></div>
                            <div class="bx--slider__filled-track"></div>
                            <div class="bx--slider__thumb" tabindex="0"></div>
                            <input id="slider" class="bx--slider__input" type="range" step="1" min="0" max="100" value="50">
                            </input>
                        </div>
                        <span class="bx--slider__range-label">100</span>
                        <input id="slider-input-box" type="number" class="bx--text-input bx--slider-text-input" placeholder="0">
                        </input>
                    </div>
                </div>
            </div>
        );
    }
}  
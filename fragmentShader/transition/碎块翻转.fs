#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

const ivec2 size = ivec2(4);
const float pause = 0.1;
const float dividerWidth = 0.05;
const vec4 bgcolor = vec4(0.0, 0.0, 0.0, 1.0);
const float randomness = 0.1;

vec4 getFromColor(vec2 pfr){
    if (imgNum==0){
        return texture2D(texture0, pfr);
    }

    if (imgNum==1){
        return texture2D(texture1, pfr);
    }

    if (imgNum==2){
        return texture2D(texture2, pfr);
    }

    if (mod(imgNum,3) == 0){
        return texture2D(texture0, pfr);
    }
    else if (mod(imgNum,3) == 1){
        return texture2D(texture1, pfr);
    }
    else if (mod(imgNum,3) == 2){
        return texture2D(texture2, pfr);
    }

    return texture2D(texture0,pfr);
}

vec4 getToColor(vec2 pto){
    if (imgNum==0){
        return texture2D(texture1, pto);
    }

    if (imgNum==1){
        return texture2D(texture2, pto);
    }

    if (imgNum==2){
        return texture2D(texture0, pto);
    }

    if (mod(imgNum,3) == 0){
        return texture2D(texture1, pto);
    }
    else if (mod(imgNum,3) == 1){
        return texture2D(texture2, pto);
    }
    else if (mod(imgNum,3) == 2){
        return texture2D(texture0, pto);
    }

    return texture2D(texture1,pto);

}

float rand (vec2 co) {
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

float getDelta(vec2 p) {
    vec2 rectanglePos = floor(vec2(size) * p);
    vec2 rectangleSize = vec2(1.0 / vec2(size).x, 1.0 / vec2(size).y);
    float top = rectangleSize.y * (rectanglePos.y + 1.0);
    float bottom = rectangleSize.y * rectanglePos.y;
    float left = rectangleSize.x * rectanglePos.x;
    float right = rectangleSize.x * (rectanglePos.x + 1.0);
    float minX = min(abs(p.x - left), abs(p.x - right));
    float minY = min(abs(p.y - top), abs(p.y - bottom));
    return min(minX, minY);
}

float getDividerSize() {
    vec2 rectangleSize = vec2(1.0 / vec2(size).x, 1.0 / vec2(size).y);
    return min(rectangleSize.x, rectangleSize.y) * dividerWidth;
}

vec4 transition(vec2 p) {
    if(time < pause) {
        float currentProg = time / pause;
        float a = 1.0;
        if(getDelta(p) < getDividerSize()) {
        a = 1.0 - currentProg;
        }
        return mix(bgcolor, getFromColor(p), a);
    } else if(time < 1.0 - pause){
        if(getDelta(p) < getDividerSize()) {
            return bgcolor;
        } else {
            float currentProg = (time - pause) / (1.0 - pause * 2.0);
            vec2 q = p;
            vec2 rectanglePos = floor(vec2(size) * q);

            float r = rand(rectanglePos) - randomness;
            float cp = smoothstep(0.0, 1.0 - r, currentProg);

            float rectangleSize = 1.0 / vec2(size).x;
            float delta = rectanglePos.x * rectangleSize;
            float offset = rectangleSize / 2.0 + delta;

            p.x = (p.x - offset)/abs(cp - 0.5)*0.5 + offset;
            vec4 a = getFromColor(p);
            vec4 b = getToColor(p);

            float s = step(abs(vec2(size).x * (q.x - delta) - 0.5), abs(cp - 0.5));
            return mix(bgcolor, mix(b, a, step(cp, 0.5)), s);
        }
    } else {
        float currentProg = (time - 1.0 + pause) / pause;
        float a = 1.0;
        if(getDelta(p) < getDividerSize()) {
            a = currentProg;
        }
        return mix(bgcolor, getToColor(p), a);
    }
}

void main(){
    FragColor = transition(uv);
}
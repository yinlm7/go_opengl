#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

const float reflection = 0.4;
const float perspective = 0.2;
const float depth = 3.0;
const vec4 black = vec4(0.0, 0.0, 0.0, 1.0);
const vec2 boundMin = vec2(0.0, 0.0);
const vec2 boundMax = vec2(1.0, 1.0);

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

bool inBounds (vec2 p) {
    return all(lessThan(boundMin, p)) && all(lessThan(p, boundMax));
}

vec2 project (vec2 p) {
    return p * vec2(1.0, -1.2) + vec2(0.0, -0.02);
}

vec4 bgColor (vec2 p, vec2 pfr, vec2 pto) {
    vec4 c = black;
    pfr = project(pfr);
    if (inBounds(pfr)) {
        c += mix(black, getFromColor(pfr), reflection * mix(1.0, 0.0, pfr.y));
    }
    pto = project(pto);
    if (inBounds(pto)) {
        c += mix(black, getToColor(pto), reflection * mix(1.0, 0.0, pto.y));
    }
    return c;
}

vec4 transition(vec2 p) {
    vec2 pfr, pto = vec2(-1.);

    float size = mix(1.0, depth, time);
    float persp = perspective * time;
    pfr = (p + vec2(-0.0, -0.5)) * vec2(size/(1.0-perspective*time), size/(1.0-size*persp*p.x)) + vec2(0.0, 0.5);

    size = mix(1.0, depth, 1.-time);
    persp = perspective * (1.-time);
    pto = (p + vec2(-1.0, -0.5)) * vec2(size/(1.0-perspective*(1.0-time)), size/(1.0-size*persp*(0.5-p.x))) + vec2(1.0, 0.5);

    if (time < 0.5) {
        if (inBounds(pfr)) {
            return getFromColor(pfr);
        }
        if (inBounds(pto)) {
            return getToColor(pto);
        }
    }

    if (inBounds(pto)) {
        return getToColor(pto);
    }

    if (inBounds(pfr)) {
        return getFromColor(pfr);
    }
    return bgColor(p, pfr, pto);
}

void main(){
    FragColor = transition(uv);
}
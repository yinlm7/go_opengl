#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

const float ratio = 1.0;
// In degrees
const float rotation = 6;
// Multiplier
const float scale = 1.2;
const float DEG2RAD = 0.03926990816987241548078304229099;

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

vec4 transition(vec2 uv) {
    // Massage parameters
    float phase = time < 0.5 ? time * 2.0 : (time - 0.5) * 2.0;
    float angleOffset = time < 0.5 ? mix(0.0, rotation * DEG2RAD, phase) : mix(-rotation * DEG2RAD, 0.0, phase);
    float newScale = time < 0.5 ? mix(1.0, scale, phase) : mix(scale, 1.0, phase);

    vec2 center = vec2(0, 0);

    // Calculate the source point
    vec2 assumedCenter = vec2(0.5, 0.5);
    vec2 p = (uv.xy - vec2(0.5, 0.5)) / newScale * vec2(ratio, 1.0);

    // This can probably be optimized (with distance())
    float angle = atan(p.y, p.x) + angleOffset;
    float dist = distance(center, p);
    p.x = cos(angle) * dist / ratio + 0.5;
    p.y = sin(angle) * dist + 0.5;
    vec4 c = time < 0.5 ? getFromColor(p) : getToColor(p);

    // Finally, apply the color
    return c + (time < 0.5 ? mix(0.0, 1.0, phase) : mix(1.0, 0.0, phase));
}

void main(){
    FragColor = transition(uv);
}
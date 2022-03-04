#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

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

vec2 offset(float time, float x, float theta) {
    float phase = time*time + time + theta;
    float shifty = 0.03*time*cos(10.0*(time+x));
    return vec2(0, shifty);
}

vec4 transition(vec2 p) {
    return mix(getFromColor(p + offset(time, p.x, 0.0)), getToColor(p + offset(1.0-time, p.x, 3.14)), time);
}

void main(){
    FragColor = transition(uv);
}
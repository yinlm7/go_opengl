#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

uniform float intensity = 0.3; // if 0.0, the image directly turn grayscale, if 0.9, the grayscale transition phase is very important

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

vec3 grayscale (vec3 color) {
return vec3(0.2126*color.r + 0.7152*color.g + 0.0722*color.b);
}

vec4 transition (vec2 uv) {
    vec4 fc = getFromColor(uv);
    vec4 tc = getToColor(uv);
    return mix(
        mix(vec4(grayscale(fc.rgb), 1.0), fc, smoothstep(1.0-intensity, 0.0, time)),
        mix(vec4(grayscale(tc.rgb), 1.0), tc, smoothstep(    intensity, 1.0, time)),
        time);
}

void main(){
    FragColor = transition(uv);
}
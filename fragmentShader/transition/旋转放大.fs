#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

in vec2 uv;
uniform float time;
uniform int imgNum;

#define PI 3.14159265359
const vec2 center = vec2(0.5, 0.5);
const float rotations = 1;
const float scale = 8;
const vec4 backColor = vec4(0.15, 0.15, 0.15, 1.0);

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


vec4 transition (vec2 uv) {

  vec2 difference = uv - center;
  vec2 dir = normalize(difference);
  float dist = length(difference);

  float angle = 2.0 * PI * rotations * time;

  float c = cos(angle);
  float s = sin(angle);

  float currentScale = mix(scale, 1.0, 2.0 * abs(time - 0.5));

  vec2 rotatedDir = vec2(dir.x  * c - dir.y * s, dir.x * s + dir.y * c);
  vec2 rotatedUv = center + rotatedDir * dist / currentScale;

  if (rotatedUv.x < 0.0 || rotatedUv.x > 1.0 ||
      rotatedUv.y < 0.0 || rotatedUv.y > 1.0)
    return backColor;

  return mix(getFromColor(rotatedUv), getToColor(rotatedUv), time);
}

void main(){
    FragColor = transition(uv);
}
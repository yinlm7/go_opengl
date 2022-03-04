#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform float time;
in vec2 uv;

float stepping(float t){
    if(t<0.)return -1.+pow(1.+t,2.);
    else return 1.-pow(1.-t,2.);
}

void main()
{
    vec2 uv1 = uv;
    vec2 fragCoord = gl_FragCoord.xy;
    vec2 iResolution = fragCoord/uv1;
    vec2 uv = (fragCoord*2.-iResolution.xy)/iResolution.y;
    FragColor = texture2D(texture0,uv1);
    uv = normalize(uv) * length(uv);
    for(int i=0;i<12;i++){
        float t = time + float(i)*3.141592/12.*(5.+1.*stepping(sin(time*3.)));
        vec2 p = vec2(cos(t),sin(t));
        p *= cos(time + float(i)*3.141592*cos(time/8.));
        vec3 col = cos(vec3(0,1,-1)*3.141592*2./3.+3.141925*(time/2.+float(i)/5.)) * 0.5 + 0.5;
        FragColor += vec4(0.05/length(uv-p*0.9)*col,1.0);
    }
    FragColor.xyz = pow(FragColor.xyz,vec3(3.));
    FragColor.w = 1.0;
}
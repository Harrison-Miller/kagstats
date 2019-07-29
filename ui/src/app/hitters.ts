// Taken directly from Base/Entities/Common/Attacks/Hitters.as
export enum Hitters {
    nothing = 0,

    //env
    crush = 1, //(required to be 1 for engine reasons)
    fall,
    water,      //just fire
    water_stun, //splash
    water_stun_force, //splash
    drown,
    fire,   //initial burst (should ignite things)
    burn,   //burn damage
    flying,

    //common actor
    stomp,
    suicide = 11, //(required to be 11 for engine reasons)

    //natural
    bite,

    //builders
    builder,

    //knight
    sword,
    shield,
    bomb,

    //archer
    stab,

    //arrows and similar projectiles
    arrow,
    bomb_arrow,
    ballista,

    //cata
    cata_stones,
    cata_boulder,
    boulder,

    //siege
    ram,

    // explosion
    explosion,
    keg,
    mine,
    mine_special,

    //traps
    spikes,

    //machinery
    saw,
    drill,

    //barbarian
    muscles,

    // scrolls
    suddengib
}


export  const HITTER_DESCRIPTION: string[] = [
    "Died",
    "Crushing",
    "Fall",
    "Drowning", //water
    "Drowning", //water_stun
    "Drowning", //water_stun_force
    "Drowning", //drown
    "Burning", //fire
    "Burning", //fire
    "Flying??",
    "Stomping",
    "Suicide",
    "Bite",
    "Pickaxe", //builder
    "Sword",
    "Shield",
    "Bomb",
    "Stab",
    "Arrow",
    "Bomb Arrow",
    "Ballista Bolt",
    "Catapult Stones",
    "Catapult Boulder",
    "Boulder",
    "Ram",
    "Explosion",
    "Keg",
    "Mine",
    "Mine", //mine_special
    "Spikes",
    "Saw",
    "Muscles",
    "Carnage"
];
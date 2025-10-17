## [EN] Description
Client-server turn-based strategy game "Troublemakers". Uses [self-writen implementation](https://github.com/Kanzu32/go-ecs) of the Enity-Component-System architectural pattern.

The game takes place between two players with a squad of several characters. Units have hit points and energy points. On his turn, the player can choose one of his units and spend any amount of energy on performing actions such as moving or attacking an opponent’s unit. Units have unique types of weapons, which determine the characteristics and type of attack.

* Dagger - increased damage to the back;
* Glaive - long-range attack in the area with the ability to hit allied units;
* Shield - increased health.

The server part performs data transfer between clients, logging, user authentication and opponent search.

## Features
* Multiple types of units;
* Settings for language, sounds and full screen mode;
* Supports local play on one device;
* User account system;
* ECS architecture.

## Technologies
* Golang;
* Data-oriented;
* Entity-Component-System;
* Ebitengine;
* Tiled;
* MongoDB.

## [RU] Описание
Клиент-серверная пошаговая стратегическая игра "Смутьяны". Использует [индивидуальную реализацию](https://github.com/Kanzu32/go-ecs) архитектурного шаблона Enity-Component-System.

Игра проходит между двумя игроками обладающими отрядом из нескольких персонажей. У юнитов есть очки жизней и энергии. В свой ход игрок может выбрать одного своего юнита и потратить любое количество энергии на совершение действий таких как передвижение или атака юнита соперника. Юниты обладают уникальными типами оружия от которого зависят характеристики и тип атаки.

* Кинжал - повышенный урон в спину;
* Глефа - дальняя атака по области с возможностью задеть своих юнитов;
* Щит - повышенное здоровье.

Серверная часть выполняет передачу данных между клиентами, логирование, аутентификацию пользователей и поиск оппонента.

## Особенности
* Несколько типов юнитов;
* Настройки языка, звуков и полноэкранного режима;
* Поддержка локальной игры на одном устройстве;
* Система аккаунтов пользователей;
* ECS архитектура.



## Screenshots
![](https://github.com/Kanzu32/strategy-game/blob/main/readme/strategy-game-1.png)
![](https://github.com/Kanzu32/strategy-game/blob/main/readme/strategy-game-2.png)
![](https://github.com/Kanzu32/strategy-game/blob/main/readme/strategy-game-3.png)
![](https://github.com/Kanzu32/strategy-game/blob/main/readme/strategy-game-4.png)

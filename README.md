# gotavern
> ** This document represents a prototype and a work in progress.**

So lately I've learned the existence of Go, and I wanted to try to test some of my development knowledge over a simple exercise written with this language.

## Creating a Purpose
The thing is that in order to create a program usually you have to find a purpose for it, not the other way around. So for the sake of this exercise, I created a fictitious scenario that I have seen repeatedly lately in some series that I've been watching.

## Adventurers Guilds
In many fantasy worlds, there are "adventurers guilds" where the person wishing to work as an adventurer needs to register in order to take the requests that the different citizens post on a mural. The quests and the adventurers have a ranking, allowing only certain adventurers to take higher quests. In these worlds, the way to take a quest is going to the register with the quest request. There the adventurer will be told who to speak with if it meets the requirements.

## Automating the Process
In this fictitious world that I have devised, the adventurer's guild requires a way to automate the whole process. They also wish to be able to portray the quests in different establishments where they have an agreement with (the tavern, the king's castle, the thieves guild, etc.). Therefore, they expect a high volume of quest requests but also the ability to scale up easily. In case they move to other cities, they would like high portability.

## Microservice Architecture
So, a Microservice architecture fits exactly their needs, and having a RabbitMQ queue to make sure none of the quest requests are lost while also allowing different consumers to deal with the high traffic of requests.

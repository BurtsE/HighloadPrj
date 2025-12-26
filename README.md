# Highload Service

Высоконагруженный сервис для обработки потоковых данных с аналитикой на основе статистических методов, развернутый в Kubernetes с автомасштабированием.

## Установка
 - Сборка образа
    - docker build -t go-service:latest .

 - Загрузка в Minikube registry
   - minikube image load go-service:latest

 - Деплой
   - helm install go-service ./deploy/helm/service --namespace metrics 

Требуется предварительная установка Redis, настройка подключения в конфигурационных файлах

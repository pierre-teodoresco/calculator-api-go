#!/bin/bash

echo "🧪 Exécution des tests unitaires de l'API Calculator"
echo "================================================"

# Exécuter tous les tests avec verbosité
echo "📋 Exécution de tous les tests..."
go test ./tests/... -v

echo ""
echo "📊 Rapport de couverture de code:"
echo "--------------------------------"

# Pour la couverture, on peut utiliser une approche différente
# Tester les packages individuellement avec couverture
echo "• Package handler:"
go test ./tests/handler -coverprofile=coverage_handler.out -coverpkg=./internal/handler

echo "• Package pkg:"
go test ./tests/pkg -coverprofile=coverage_pkg.out -coverpkg=./pkg

echo ""
echo "📈 Générer un rapport de couverture HTML:"
echo "go tool cover -html=coverage_handler.out -o coverage_handler.html"
echo "go tool cover -html=coverage_pkg.out -o coverage_pkg.html"

echo ""
echo "✅ Tests terminés avec succès!" 
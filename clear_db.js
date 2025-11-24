// Clear the searchengine database
use searchengine;

print("=== CLEARING SEARCH ENGINE DATABASE ===");
print("Pages before deletion:", db.pages.countDocuments());

// Delete all pages
var result = db.pages.deleteMany({});
print("Deleted", result.deletedCount, "pages");

// Delete all index entries  
var indexResult = db.index.deleteMany({});
print("Deleted", indexResult.deletedCount, "index entries");

// Delete all pagerank entries
var prResult = db.pagerank.deleteMany({});
print("Deleted", prResult.deletedCount, "pagerank entries");

print("Pages after deletion:", db.pages.countDocuments());
print("=== DATABASE CLEARED SUCCESSFULLY ===");